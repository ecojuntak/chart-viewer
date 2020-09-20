<template>
  <v-container>
    <two-versions-selector 
      @firstTemplatesChanged="setFirstTemplates"
      @secondTemplatesChanged="setSecondTemplates"
    />

    <v-row v-if="firstTemplates.length != 0 && secondTemplates.length != 0">
      <diff-viewer :firstTemplates="firstTemplates" :secondTemplates="secondTemplates"> </diff-viewer>
    </v-row>

  </v-container>
</template>

<script>
  import twoVersionsSelector from '../components/TwoVersionsSelector'
  import diffViewer from '../components/DiffViewer'

  export default {
    name: 'Compare',
    components: {
      diffViewer,
      twoVersionsSelector
    },
    data () {
      return {
        firstTemplates: [],
        secondTemplates: [],
      }
    },
    methods: {
      setFirstTemplates(templates) {
        this.firstTemplates = templates
      },
      setSecondTemplates(templates) {
        this.secondTemplates = templates
      }
    }  
  }
</script>