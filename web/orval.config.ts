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
        schemas: {
          'db.CourseWithTranslation': {
            required: [
              'uuid',
              'created_at',
              'modified_at',
              'subject',
              'price',
              'discount',
              'is_active',
              'difficulty',
              'translation',
            ],
          },
        },
      },
      prettier: true,
    },
  },
});
