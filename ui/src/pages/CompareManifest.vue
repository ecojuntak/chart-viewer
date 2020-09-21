<template>
  <v-container>
    <two-versions-selector
      @selectedRepoChanged="setRepo"
      @selectedChartChanged="setChart"
      @firstValuesChanged="setFirstValues"
      @secondValuesChanged="setSecondValues"
      @firstSelectedVersionChanged="setFirstVersion"
      @secondSelectedVersionChanged="setSecondVersion"
    />

    <v-row v-if="progressing">
      <v-progress-linear
        indeterminate
        color="blue darken-2"
      ></v-progress-linear>
    </v-row>

    <v-row v-if="firstManifests.length == 0 && secondManifests.length == 0">
      <v-col v-if="firstValues != ''" cols="6">
        <v-alert type="error" dense outlined cols="12" v-if="firstErrorMessage != ''">
          {{ firstErrorMessage }}
        </v-alert>
        <code>values.yaml</code> for version {{ firstVersion }}
        <prism-editor
          class="my-editor overflow-x-auto" 
          v-model="firstValues"
          :highlight="highlighter"
          line-numbers
          :readonly="false"
        />
      </v-col>

      <v-col v-if="secondValues != ''" cols="6">
        <v-alert type="error" dense outlined cols="12" v-if="secondErrorMessage != ''">
          {{ secondErrorMessage }}
        </v-alert>
        <code>values.yaml</code> for version {{ secondVersion }}
        <prism-editor
          class="my-editor overflow-x-auto" 
          v-model="secondValues"
          :highlight="highlighter"
          line-numbers
          :readonly="false"
        />
      </v-col>

      <v-btn block color="primary" dark 
        @click="renderBothManifest()" class="mt-2"
        v-if="secondValues != '' && firstValues != ''">
        Compare Manifest
      </v-btn>
    </v-row>
    <v-row v-if="firstManifests.length != 0 && secondManifests.length != 0">
      <diff-viewer :firstTemplates="firstManifests" :secondTemplates="secondManifests"> </diff-viewer>
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

  import api from '../api/api'
  import twoVersionsSelector from '../components/TwoVersionsSelector'
  import diffViewer from '../components/DiffViewer'
  
  export default {
    name: 'CompareManifest',
    components: {
      twoVersionsSelector,
      diffViewer,
      PrismEditor
    },
    data () {
      return {
        repo: "",
        chart: "",
        firstValues: "",
        secondValues: "",
        firstVersion: "",
        secondVersion: "",
        firstManifests: [],
        secondManifests: [],
        progressing: false,
        firstErrorMessage: "",
        secondErrorMessage: "",
      }
    },
    methods: {
      async renderBothManifest() {
        this.progressing = true

        const firstResponse = await this.renderManifests(this.firstValues, this.firstVersion)
        if(firstResponse.status == 500) {
          this.firstErrorMessage = firstResponse.data.error
        } else {
          this.firstManifests = firstResponse.data.manifests
        }

        const secondResponse = await this.renderManifests(this.secondValues, this.secondVersion)
        if(secondResponse.status == 500) {
          this.secondErrorMessage = secondResponse.data.error
        } else {
          this.secondManifests = secondResponse.data.manifests
        }

        this.progressing = false
      },
      async renderManifests (values, version) {
        const escaptedValues = escape(values)
        return await api.renderManifest(this.repo, this.chart, version, escaptedValues)
      },
      setFirstValues(values) {
        this.firstValues = values
      },
      setSecondValues(values) {
        this.secondValues = values
      },
      setFirstVersion(version) {
        this.firstVersion = version
        this.firstManifests = []
        this.secondManifests = []
      },
      setSecondVersion(version) {
        this.secondVersion = version
        this.firstManifests = []
        this.secondManifests = []
      },
      setRepo(repo) {
        this.repo = repo
        this.firstValues = ""
        this.secondValues = ""
        this.firstManifests = []
        this.secondManifests = []
      },
      setChart(chart) {
        this.chart = chart
        this.firstValues = ""
        this.secondValues = ""
        this.firstManifests = []
        this.secondManifests = []
      },
      highlighter(values) {
        return highlight(values, languages.yaml);
      },
    }
  }
</script>

<style scoped>
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