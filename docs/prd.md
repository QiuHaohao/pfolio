# PRD - pfolio

## Background

Robo-advisor products like Syfe Wealth provides functionalities like deposit, withdrawal and rebalancing with a range of portfolio to choose from, and charges a percentage of value of assets managed as fee(0.4% for accounts with more than SGD 100k managed). Under this pricing model, while fee increases linearly as value of assets managed increases, the value that the produce provides does not really increase together with the fee - after all, all the product does is just buying, holding and selling assets for the client.

With a portfolio of US$200,000, Syfe Wealth can charge US$800 annually, which is more than US$65 monthly. Assuming the portfolio consists of 20 index funds, and deposit is made monthly, the same portfolio can be maintained with less than US$10 monthly commission, or even free of charge on TDAmeritrade. Although working directly with brokers introduces fund transfer fees and other costs, it still is generally a cheaper and more transparent option than most robo-advisors, especially for larger portfolios.

## Objective

The purpose of this project is to build an alternative to robo-advisors that directly interacts with brokers to view portfolio information and send orders. The first version of this product will be a command-line tool. It can potentially be later extended to a service.

## Functional Requirements

- Model
    - Add a new model
        - Static model
            - List of symbol - weight
            - equivalent instruments
        - Dynamic model
            - Derived from external source
            - synchronize
            - e.g. Syfe portfolios are publicly accessible
    - List models
- Portfolio
    - Add a new portfolio
        - model
        - currency
        - broker
    - Remove a portfolio
    - List portfolios
    - View a portfolio
    - Shadow portfolio
        - one that uses broker as data source, but executed by updating a file
- Orders
    - buy/sell
    - support preview
    - User confirmation with live update?
- Configuration
    - editor
    - broker
    - buy/sell strategy

## Non-functional Requirements

- The ability to integrate with multiple brokers
    - Abstract the broker interface away

## Interface

### Model

#### Create

```bash
> pfolio model create -n <model_name>
# Create using template
# Open editor
# Validate
```

#### Update

```bash
> pfolio model update <model_name>
# Open editor
# Validate
```

#### View

```bash
> pfolio model view <model_name>
# Open with readonly
```

#### Remove

```bash
> pfolio model rm <model_name>
# Remove entity file
```

#### List

```bash
> pfolio model list
# Show a list of models
```

### Metamodel

#### Create

```bash
> pfolio metamodel create syfe --type <type> <metamodel_name>
# Create a meta model 
# Sync model
```

#### Update

```bash
> pfolio metamodel update <metamodel_name>
# Open editor
# Validate
```

#### View

```bash
> pfolio metamodel view <metamodel_name>
# Open with readonly
```

#### Remove

```bash
> pfolio metamodel rm <metamodel_name>
# Remove entity file
```

#### List

```bash
> pfolio metamodel list
# Show a list of metamodels
```

#### Sync

```bash
> pfolio metamodel sync <metamodel>
# Update the model corresponding to a metamodelUpdate the model corresponding to a metamodel
```

### Portfolio

#### Create

```bash
> pfolio portfolio create <portfolio_name>
# Create portfolio file with template
# Open Editor
# Validate
```

#### Update

```bash
> pfolio portfolio edit <portfolio_name>
# Open editor
# Validate
```

#### View

```bash
> pfolio portfolio view <portfolio_name>
# Open with readonly
```

#### Remove

```bash
> pfolio portfolio rm <portfolio_name>
# Remove portfolio file
```

#### List

```bash
> pfolio portfolio list
# Show a list of portfolios
```

### Orders

#### Buy

```bash
> pfolio buy --execute --amount <amount> --all <portfolio_name>
```

### Sell

```bash
> pfolio sell --execute --amount <amount> --percentage <percentage> --all <portfolio_name>
```

### Configure

```bash
> pfolio config
# Open editor
# Validate
```
