const axios = require('axios')

const baseURL = process.env.VUE_APP_API_SERVER_HOST

async function fetchRepos() {
    return await axios.get(baseURL + '/repos')
}

async function fetchCharts(repoName) {
    return await axios.get(baseURL + '/charts/' + repoName)
}

async function fetchChart(repoName, chartName, chartVersion) {
    return await axios.get(baseURL + '/charts/' + repoName + '/' + chartName + '/' + chartVersion)
}

async function renderManifest(repoName, chartName, chartVersion, values) {
    return await axios.get(baseURL + '/charts/manifests/render/' + repoName + '/' + chartName + '/' + chartVersion + '?values=' + values)
    .then(res => res)
    .catch(err => err.response)
}

module.exports = {
    fetchRepos,
    fetchCharts,
    fetchChart,
    renderManifest
}
