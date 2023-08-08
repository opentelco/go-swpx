
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: '../graph/src/*.graphqls',
  generates: {
    'src/gql/': {
      preset: 'client',
      plugins: [],
      config: {
        useTypeImports: true
      }
    },
    './graphql.schema.json': {
      plugins: ['introspection']
    }
  }
};

export default config;
