<template>
  <v-container>
    <one-version-selector-with-kube-version
      @templatesChanged="setTemplate"
    />
    
    <v-row v-if="templates.length !== 0">
      <chart-viewer :templates="templates" />
    </v-row>
    
  </v-container>
</template>

<script>

  import chartViewer from '../components/ChartViewer'
  import oneVersionSelectorWithKubeVersion from '../components/OneVersionSelectorWithKubeVersion'

  export default {
    name: 'Home',
    components: {
      chartViewer,
      oneVersionSelectorWithKubeVersion,
    },
    data () {
      return {
        templates: [],
      }
    },
    methods: {
      setTemplate(templates) {
        this.templates = templates
      }
    }
  }
</script>