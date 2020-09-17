<template>
  <v-container>
    <v-row align="center">
      <v-col class="d-flex" cols="4" sm="4">
        <v-autocomplete
          :items="repos"
          label="Repo Name"
          :item-text="contructChartName"
          item-value="name"
          v-model="selectedRepo"
          @change="fetchChartList"
        ></v-autocomplete>
      </v-col>

      <v-col class="d-flex" cols="4" sm="4">
        <v-autocomplete
          :items="charts"
          label="Chart Name"
          item-text="name"
          v-model="selectedChart"
          @change="fetchVersionList"
        ></v-autocomplete>
      </v-col>

      <v-col class="d-flex" cols="4" sm="4">
        <v-autocomplete
          :items="versions"
          label="Chart Version"
          v-model="selectedVersion"
          @change="fetchChartDetail"
        ></v-autocomplete>
      </v-col>
    </v-row>

    <v-row v-if="values !='' && manifests.length == 0">
      <v-col cols="8">
        <p>
          Costomize the <code>values.yaml</code> below
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
        {{ generatedUrl }}
        <v-btn icon @click="copyGeneratedURL()">
          <v-icon> mdi-content-copy</v-icon>
        </v-btn>
        <v-alert
            v-if="copied"
            dense
            rounded
            text
            type="success"
          >
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
  import yaml from 'json-to-pretty-yaml'
  import chartViewer from '../components/ChartViewer'

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
      chartViewer
    },
    data () {
      return {
        repos: [],
        selectedRepo: "",
        charts: [],
        selectedChart: "",
        versions: [],
        selectedVersion: "",
        values: "",
        progressing: false,
        manifests: [],
        generatedUrl: "",
        copied: false,
      }
    },
    mounted() {
      this.fetchRepoList()
    },
    methods: {
      async fetchRepoList() {
        this.resetState()
        const response = await api.fetchRepoList()
        this.repos = response.data
      },
      async fetchChartList() {
        this.resetState()
        const response = await api.fetchChartList(this.selectedRepo)
        this.charts = response.data
      },
      fetchVersionList() {
        this.resetState()
        for(let i=0; i < this.charts.length; i++) {
          if(this.charts[i].name === this.selectedChart) {
            this.versions = this.charts[i].versions
            break
          }
        }
      },
      async fetchChartDetail() {
        this.values = ""
        this.progressing = true
        const response = await api.fetchChartDetail(this.selectedRepo, this.selectedChart, this.selectedVersion)
        this.progressing = false

        this.values = yaml.stringify(response.data.values)
      },
      async getManifest(){
        const test = escape(this.values)
     
        this.progressing = true
        const response = await api.fetchManifest(this.selectedRepo, this.selectedChart, this.selectedVersion, test)
        this.progressing = false
        this.manifests = response.data.manifests
        this.generatedUrl = process.env.VUE_APP_API_SERVER_HOST + "/manifest/" + response.data.url
      },
      contructChartName(repo) {
        return repo.name + " (" + repo.url + ")"
      },
      resetState() {
        this.versions = []
        this.manifests = []
        this.copied = false
        this.generatedUrl = ""
        this.values = ""
      },
      highlighter(values) {
        return highlight(values, languages.yaml);
      },
      copyGeneratedURL() {
        navigator.clipboard.writeText(this.generatedUrl);
        this.copied = true
      }
    }
  }
</script>