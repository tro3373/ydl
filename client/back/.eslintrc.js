module.exports = {
  root: true,
  env: {
    es2020: true,
    node: true,
    browser: true,
  },
  plugins: ['vue', 'jest'],
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
    parser: 'babel-eslint',
  },
  extends: [
    'eslint:recommended',
    'plugin:jest/recommended',
    'plugin:vue/essential',
    'prettier',
    // "@vue/prettier",
    // "plugin:vue/base"
  ],
  overrides: [
    {
      files: ['**/__tests__/*.{j,t}s?(x)', '**/test/*.spec.{j,t}s?(x)'],
      env: {
        'jest/globals': true,
      },
    },
  ],
  // rules: {
  //   // 0は許可,1は警告,2はエラー
  //   'vue/html-closing-bracket-newline': [
  //     'error',
  //     {
  //       multiline: 'never'
  //     }
  //   ],
  //   // 不要なカッコは消す
  //   'no-extra-parens': 1,
  //   // 無駄なスペースは削除
  //   'no-multi-spaces': 2,
  //   // 不要な空白行は削除。2行開けてたらエラー
  //   'no-multiple-empty-lines': [
  //     2,
  //     {
  //       max: 1
  //     }
  //   ],
  //   // 関数とカッコはあけない(function hoge() {/** */})
  //   'func-call-spacing': [2, 'never'],
  //   // true/falseを無駄に使うな
  //   'no-unneeded-ternary': 2,
  //   // セミコロン強制
  //   semi: [2, 'always'],
  //   // 文字列はシングルクオートのみ
  //   quotes: [2, 'single'],
  //   // varは禁止
  //   'no-var': 2,
  //   // jsのインデントは２
  //   indent: [2, 2],
  //   // かっこの中はスペースなし！違和感
  //   'space-in-parens': [2, 'never'],
  //   // コンソールは許可
  //   'no-console': 0,
  //   // カンマの前後にスペース入れる？
  //   'comma-spacing': 2,
  //   // 配列のindexには空白入れるな(hogehoge[ x ])
  //   'computed-property-spacing': 2,
  //   // キー
  //   'key-spacing': 2,
  //   // キーワードの前後には適切なスペースを
  //   'keyword-spacing': 2,
  //   // line separetor: lf
  //   'linebreak-style': [2, 'unix']
  // }
};
