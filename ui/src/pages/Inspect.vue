<template>
  <v-container>
    <one-version-selector 
      @templatesChanged="setTemplate"
    />
    
    <v-row v-if="templates.length != 0">
      <chart-viewer :templates="templates" />
    </v-row>
    
  </v-container>
</template>

<script>

  import chartViewer from '../components/ChartViewer'
  import oneVersionSelector from '../components/OneVersionSelector'

  export default {
    name: 'Home',
    components: {
      chartViewer,
      oneVersionSelector,
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