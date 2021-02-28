const axios = require('axios')

async function fetchRepos() {
    return await axios.get('/api/v1/repos')
}

async function fetchCharts(repoName) {
    return await axios.get('/api/v1/charts/' + repoName)
}

async function fetchChart(repoName, chartName, chartVersion) {
    return await axios.get('/api/v1/charts/' + repoName + '/' + chartName + '/' + chartVersion)
}

async function renderManifest(repoName, chartName, chartVersion, values) {
    return await axios.post('/api/v1/charts/manifests/render/' + repoName + '/' + chartName + '/' + chartVersion, {'values': values})
    .then(res => res)
    .catch(err => err.response)
}

module.exports = {
    fetchRepos,
    fetchCharts,
    fetchChart,
    renderManifest
}
