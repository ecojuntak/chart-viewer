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
          @change="fetchChart"
        ></v-autocomplete>
      </v-col>
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

  export default {
    name: 'OneVersionSelector',
    data () {
      return {
        repos: [],
        selectedRepo: this.$route.query.repo || '',
        charts: [],
        selectedChart: this.$route.query.chart || '',
        versions: [],
        selectedVersion: this.$route.query.version || '',
        values: '',
        templates: [],
        progressing: false,
        initailLoad: true
      }
    },
    async mounted() {
      await this.fetchRepoList()
      if(this.selectedRepo !== '') {
        await this.fetchChartList()
      }

      if(this.selectedRepo !== '' && this.selectedChart !== '') {
        await this.fetchVersionList()
      }

      if(this.selectedRepo !== '' && this.selectedChart !== '' && this.selectedVersion !== '') {
        await this.fetchChart()
      }

      this.initailLoad = false
    },
    methods: {
      async fetchRepoList() {
        if(!this.initailLoad) {
          this.resetState()
        }

        const response = await api.fetchRepos()
        this.repos = response.data
      },
      async fetchChartList() {
        if(!this.initailLoad) {
          this.resetState()
          this.selectedChart = ""
          this.selectedVersion = ""
        }

        const response = await api.fetchCharts(this.selectedRepo)
        this.charts = response.data
        this.updateQueryParams()
      },
      async fetchChart() {
        this.values = ""
        this.templates = []

        this.progressing = true
        const response = await api.fetchChart(this.selectedRepo, this.selectedChart, this.selectedVersion)
        this.progressing = false
        
        this.values = yaml.stringify(response.data.values)
        const templates = response.data.templates
        this.templates = this.simplifyTemplateName(templates)
        this.templates.push({
          name: "values.yaml",
          content: this.values
        })
        this.updateQueryParams()
      },
      fetchVersionList() {
        this.templates = []
        for(let i=0; i < this.charts.length; i++) {
          if(this.charts[i].name === this.selectedChart) {
            this.versions = this.charts[i].versions
            break
          }
        }
        this.updateQueryParams()
      },
      simplifyTemplateName(templates) {
        var temps = []
        templates.forEach((template) => {
          const newName = template.name.replace("templates/", "")
          temps.push({
            name: newName,
            content: template.content
          })
        })

        return temps
      },
      contructChartName(repo) {
        return repo.name + " (" + repo.url + ")"
      },
      resetState() {
        this.versions = []
        this.templates = []
      },
      updateQueryParams() {
        this.$router.push(
          {
            path: this.$route.path,
            query: {
              repo: this.selectedRepo,
              chart: this.selectedChart,
              version: this.selectedVersion
            } 
          }
        )
      }
    },
    watch: {
      templates() {
        this.$emit("templatesChanged", this.templates);
      },
      values() {
        this.$emit("valuesChanged", this.values)
      },
      selectedRepo(){
        this.$emit("repoChanged", this.selectedRepo)
      },
      selectedChart(){
        this.$emit("chartChanged", this.selectedChart)
      },
      selectedVersion(){
        this.$emit("versionChanged", this.selectedVersion)
      }
    }
  }
</script>