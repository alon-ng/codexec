import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: {
      target: '../api/swagger.json',
    },
    output: {
      mode: 'tags-split',
      target: 'app/api/generated/endpoints.ts',
      schemas: 'app/api/generated/model',
      client: 'react-query',
      mock: true,
      override: {
        mutator: {
          path: 'app/lib/axios.ts',
          name: 'customInstance',
        },
        operations: {},
      },
      prettier: true,
    },
  },
});
