/** Root ESLint configuration shared by all JS/TS workspaces. */
module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  extends: ['eslint:recommended', 'plugin:@typescript-eslint/recommended'],
  env: { node: true, es2022: true, browser: true },
  ignorePatterns: ['dist', '.next', 'node_modules', 'coverage'],
  rules: {
    'no-console': ['warn', { allow: ['warn', 'error'] }],
    eqeqeq: ['error', 'always'],
    '@typescript-eslint/no-explicit-any': 'error',
    '@typescript-eslint/consistent-type-imports': 'warn',
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
  },
  overrides: [
    {
      files: ['scripts/**/*.js', 'apps/worker/**/*.js', '*.config.js'],
      rules: { 'no-console': 'off', '@typescript-eslint/no-var-requires': 'off' },
    },
  ],
};
