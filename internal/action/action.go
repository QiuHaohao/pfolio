package action

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/qiuhaohao/pfolio/internal/cli"
	"github.com/qiuhaohao/pfolio/internal/clock"
	"github.com/qiuhaohao/pfolio/internal/config"
	"github.com/qiuhaohao/pfolio/internal/db"
	"github.com/qiuhaohao/pfolio/internal/editor"
	"github.com/qiuhaohao/pfolio/internal/metamodel"
	"github.com/qiuhaohao/pfolio/internal/table"
	"github.com/qiuhaohao/pfolio/internal/view"
)

var (
	ErrModelNotFound           = errors.New("model not found")
	ErrDuplicatedModelName     = errors.New("model name already exists")
	ErrDuplicatedMetamodelName = errors.New("metamodel name already exists")
)

type Action struct {
	db     db.Database
	editor editor.RetryEditor
	stdin  *os.File
	stdout *os.File
	stderr *os.File
}

func New(db db.Database, editor editor.RetryEditor, stdin *os.File, stdout *os.File, stderr *os.File) Action {
	return Action{
		db:     db,
		editor: editor,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func NewDefaultAction() Action {
	return Action{
		db: db.DefaultDB(),
		editor: editor.NewYamlObjEditorWithRetry(
			editor.WithEditor(viper.GetString(config.KeyEditor)),
			editor.WithStdin(os.Stdin),
			editor.WithStdout(os.Stdout),
			editor.WithStderr(os.Stderr)),
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (a Action) states() db.States {
	return a.db.States()
}

func (a Action) CreateModel(name string, initialView view.ModelEditView) (err error) {
	if _, exists := a.states().GetModel(name); exists {
		return ErrDuplicatedModelName
	}

	obj, err := a.editor.EditWithRetry(
		initialView,
		func(v interface{}) error {
			return v.(view.ModelEditView).Validate()
		})
	if err != nil {
		return
	}

	v := obj.(view.ModelEditView)

	a.states().SetModel(name, newModel(v, false))

	if err = a.db.Save(); err != nil {
		return
	}

	_, err = fmt.Fprintf(a.stdout, "Model %s successfully created!\n", cli.Highlight(name))

	return
}

func (a Action) UpdateModel(name string) (err error) {
	m, ok := a.states().GetModel(name)
	if !ok {
		return ErrModelNotFound
	}

	obj, err := a.editor.EditWithRetry(
		view.NewModelEditViewFromDB(m),
		func(v interface{}) error {
			return v.(view.ModelEditView).Validate()
		})
	if err != nil {
		return
	}

	v := obj.(view.ModelEditView)

	m.Entries = db.ModelEntries(v)
	m.UpdateTime = clock.Now()

	a.states().SetModel(name, m)

	if err = a.db.Save(); err != nil {
		return
	}

	_, err = fmt.Fprintf(a.stdout, "Model %s successfully updated!\n", cli.Highlight(name))

	return
}

func (a Action) ViewModel(name string) error {
	m, ok := a.states().GetModel(name)
	if !ok {
		return ErrModelNotFound
	}

	cli.PrintDivider(a.stdin)
	fmt.Printf("Model Name: %s", cli.Highlight(name))
	if m.IsDerivedFromMetaModel {
		fmt.Printf("(Derived from metamodel)")
	}
	fmt.Println()
	fmt.Printf("Create Time: %s\n", m.CreateTime.Format(time.RFC822))
	fmt.Printf("Update Time: %s\n", m.UpdateTime.Format(time.RFC822))

	cli.PrintDivider(a.stdin)
	tbl := table.New(a.stdin, "Instrument", "Weight", "Percentage", "Equivalents")

	sort.Slice(m.Entries, func(i, j int) bool {
		return m.Entries[i].Weight > m.Entries[j].Weight ||
			(m.Entries[i].Weight == m.Entries[j].Weight &&
				m.Entries[i].InstrumentIdentifier < m.Entries[i].InstrumentIdentifier)
	})

	totalWeight := m.Entries.TotalWeight()

	for _, e := range m.Entries {
		tbl.AddRow(
			e.InstrumentIdentifier,
			fmt.Sprintf("%.2f", e.Weight),
			fmt.Sprintf("%.2f%%", 100*e.Weight/totalWeight),
			strings.Join(e.GetStringsEquivalentInstruments(), ", "))
	}

	tbl.Print()

	return nil
}

type ModelSortOrder int

const (
	ModelSortOrderByName ModelSortOrder = iota
	ModelSortOrderByCreateTime
	ModelSortOrderByUpdateTime
)

func (a Action) ListModels(sortOrder ModelSortOrder, descending bool) {
	type modelWithName struct {
		name  string
		model db.Model
	}

	cli.PrintDivider(a.stdin)
	tbl := table.New(a.stdin, "Model Name", "Create Time", "Update Time")

	modelsWithName := make([]modelWithName, 0)
	for name, m := range a.states().GetModels() {
		modelsWithName = append(modelsWithName, modelWithName{name: name, model: m})
	}
	var sortFn func(i, j int) bool

	switch sortOrder {
	case ModelSortOrderByCreateTime:
		sortFn = func(i, j int) bool {
			return modelsWithName[i].model.CreateTime.Before(modelsWithName[j].model.CreateTime)
		}
	case ModelSortOrderByUpdateTime:
		sortFn = func(i, j int) bool {
			return modelsWithName[i].model.UpdateTime.Before(modelsWithName[j].model.UpdateTime)
		}
	default:
		sortFn = func(i, j int) bool {
			return modelsWithName[i].name < (modelsWithName[j].name)
		}
	}

	if descending {
		ascSortFn := sortFn
		sortFn = func(i, j int) bool {
			return !ascSortFn(i, j)
		}
	}

	sort.Slice(modelsWithName, sortFn)

	for _, m := range modelsWithName {
		tbl.AddRow(m.name, m.model.CreateTime.Format(time.RFC822), m.model.UpdateTime.Format(time.RFC822))
	}

	tbl.Print()

	return
}

func (a Action) RemoveModels(names ...string) error {
	for _, name := range names {
		if _, ok := a.states().GetModel(name); !ok {
			fmt.Printf("Model %s does not exist.\n", cli.Highlight(name))
			continue
		}
		a.states().RemoveModel(name)
		fmt.Printf("Model %s removed.\n", cli.Highlight(name))
	}

	return a.db.Save()
}

func (a Action) CreateSyfeMetamodel(name string, initialView view.SyfeMetamodelEditView) (err error) {
	if _, exists := a.states().GetModel(name); exists {
		return ErrDuplicatedModelName
	}

	if _, exists := a.states().GetMetamodel(name); exists {
		return ErrDuplicatedMetamodelName
	}

	var entries db.ModelEntries

	obj, err := a.editor.EditWithRetry(
		initialView,
		func(obj interface{}) error {
			v := obj.(view.SyfeMetamodelEditView)
			if err = v.Validate(); err != nil {
				return err
			}
			mm := metamodel.NewSyfeMetamodel(v.SyfeMetamodelInfo, v.CommonMetamodelInfo)
			if entries, err = mm.GetModelEntries(); err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return
	}

	v := obj.(view.SyfeMetamodelEditView)

	a.states().SetModel(name, newModel(entries, true))
	a.states().SetMetamodel(name, newSyfeMetamodel(v.SyfeMetamodelInfo, v.CommonMetamodelInfo))

	if err = a.db.Save(); err != nil {
		return
	}

	_, err = fmt.Fprintf(a.stdout, "Syfe metamodel %s successfully created!\n", cli.Highlight(name))

	return
}

func newModel(entries []db.ModelEntry, isDerivedFromMetaModel bool) db.Model {
	return db.Model{
		Entries:                entries,
		IsDerivedFromMetaModel: isDerivedFromMetaModel,
		CreateTime:             clock.Now(),
		UpdateTime:             clock.Now(),
	}
}

func newSyfeMetamodel(syfeInfo db.SyfeMetamodelInfo, commonInfo db.CommonMetamodelInfo) db.Metamodel {
	return db.Metamodel{
		MetamodelType:       metamodel.TypeSyfe,
		CommonMetamodelInfo: commonInfo,
		SyfeMetamodelInfo:   syfeInfo,
		CreateTime:          clock.Now(),
		UpdateTime:          clock.Now(),
	}
}
