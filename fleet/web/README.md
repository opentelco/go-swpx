# OpenTelco CRM (crm)

A OpenTelco CRM system

GRAPHQL_URI=https://prod.example.com/graphql quasar build
GRAPHQL_URI=https://localhost:1336/query quasar dev

generate the gql TS:
`npx graphql-code-generator`

## Install the dependencies
```bash
yarn
# or
npm install
```

### Start the app in development mode (hot-code reloading, error reporting, etc.)
```bash
quasar dev
```


### Lint the files
```bash
yarn lint
# or
npm run lint
```


### Format the files
```bash
yarn format
# or
npm run format
```



### Build the app for production
```bash
quasar build
```

### Customize the configuration
See [Configuring quasar.config.js](https://v2.quasar.dev/quasar-cli-webpack/quasar-config-js).
