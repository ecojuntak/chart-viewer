<template>
  <v-container>
    <one-version-selector
      @valuesChanged="setValues"
      @repoChanged="setRepo"
      @chartChanged="setChart"
      @versionChanged="setVersion"
    />

    <v-row v-if="values !='' && manifests.length == 0">
      <v-alert type="error" dense outlined cols="12" v-if="errorMessage != ''">
        {{ errorMessage }}
      </v-alert>
      <v-col cols="12">
        <p>
          Customize the <code>values.yaml</code> below
        </p>
        <prism-editor
          class="my-editor overflow-x-auto" 
          v-model="values"
          :highlight="highlighter"
          line-numbers
          :readonly="false"
        />
        <v-btn block color="primary" dark @click="getManifest()" class="mt-2">Render Manifest</v-btn>
      </v-col>
    </v-row>

    <v-row v-if="manifests.length != 0">
      <div class="d-block pa-2">
        <v-btn rounded color="primary" dark small @click="manifests = []">
          <v-icon left>mdi-pencil</v-icon> Edit values.yaml
        </v-btn>
      </div>
      <div class="d-flex pa-2">
        <code>{{ generatedCommand }}</code>
        <v-btn icon @click="copyGeneratedURL()">
          <v-icon> mdi-content-copy</v-icon>
        </v-btn>
        <v-alert v-if="copied" dense rounded text type="success">
          URL copied
        </v-alert>
      </div>
      <chart-viewer :templates="manifests"> </chart-viewer>
    </v-row>

    <v-row v-if="progressing">
      <v-progress-linear
        indeterminate
        color="blue darken-2"
      ></v-progress-linear>
    </v-row>

  </v-container>
</template>

<script>
  import api from '../api/api'
  import chartViewer from '../components/ChartViewer'
  import oneVersionSelector from '../components/OneVersionSelector'

  import { PrismEditor } from 'vue-prism-editor'
  import 'vue-prism-editor/dist/prismeditor.min.css'

  import { highlight, languages } from 'prismjs/components/prism-core'
  import 'prismjs/components/prism-clike'
  import 'prismjs/components/prism-yaml'
  import 'prismjs/themes/prism-tomorrow.css'

  export default {
    name: 'Render',
    components: {
      PrismEditor,
      chartViewer,
      oneVersionSelector
    },
    data () {
      return {
        repo: "",
        chart: "",
        version: "",
        values: "",
        progressing: false,
        manifests: [],
        generatedCommand: "",
        copied: false,
        errorMessage: ""
      }
    },
    methods: {
      setValues(values) {
        this.values = values
      },
      setRepo(repo) {
        this.repo = repo
        this.resetState()
      },
      setChart(chart) {
        this.chart = chart
        this.resetState()
      },
      setVersion(version) {
        this.version = version
        this.resetState()
      },
      async getManifest(){
        const values = escape(this.values)
     
        this.progressing = true
        const response = await api.renderManifest(this.repo, this.chart, this.version, values)
        this.progressing = false

        if(response.status == 500) {
          this.errorMessage = response.data.error
        } else {
          this.manifests = response.data.manifests
          this.generatedCommand = "kubectl apply -f " + process.env.VUE_APP_API_SERVER_HOST + response.data.url
        }
      },
      contructChartName(repo) {
        return repo.name + " (" + repo.url + ")"
      },
      resetState() {
        this.versions = []
        this.manifests = []
        this.copied = false
        this.generatedCommand = ""
        this.values = ""
        this.errorMessage = ""
      },
      highlighter(values) {
        return highlight(values, languages.yaml);
      },
      copyGeneratedURL() {
        navigator.clipboard.writeText(this.generatedCommand);
        this.copied = true
      }
    }
  }
</script>