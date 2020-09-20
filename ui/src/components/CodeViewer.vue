<template>
  <v-container>
    {{ message }}
    <v-row>
      <v-col>
        <prism-editor
          class="my-editor overflow-x-auto" 
          v-model="code"
          line-numbers
          :highlight="highlighter"
          :readonly="readonly"
        >
        </prism-editor>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  import { PrismEditor } from 'vue-prism-editor';
  import 'vue-prism-editor/dist/prismeditor.min.css';
  import { highlight, languages } from 'prismjs/components/prism-core';
  import 'prismjs/components/prism-clike';
  import 'prismjs/components/prism-yaml';
  import 'prismjs/themes/prism-tomorrow.css';
  
  export default {
    name: "CodoViewer",
    components: {
      PrismEditor
    },
    data() {
      return {
        data: this.code
      }
    },
    props: [
      "code",
      "readonly",
      "message"
    ],
    methods: {
      highlighter(values) {
        return highlight(values, languages.yaml);
      },
    },
  }
</script>

<style>
  .my-editor {
    background: #2d2d2d;
    color: #ccc;

    font-family: Fira code, Fira Mono, Consolas, Menlo, Courier, monospace;
    font-size: 16px;
    line-height: 1.5;
    padding: 5px;

    max-height: 1000px;
  }
</style>