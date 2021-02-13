package assets

//HomeHTML contains the bootstrapping HTML for Swaggerbond
var HomeHTML = `
<!DOCTYPE html>
<html>
  <head>
    <title>Swaggerbond</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="/assets/bootstrap.css" rel="stylesheet">
    <script defer src="/assets/vue.js"> </script>
    <script defer src="/assets/swaggerbond.js"> </script>
  </head>
  <body>
    <div id="app" class="container-fluid">
      <div class="row pt-3 justify-content-center">
          <div class="col-12 col-sm-10 col-md-8 col-lg-6 col-xl-5">
            <search-input :searching="searching" @term-entered="searchServices($event)"></search-input>
          </div>
      </div>   
      <div class="row pt-3">
          <div class="col-12">
            <results-view :results="results"></results-view>
          </div>
      </div>
    </div>  
  </body>
</html>`

//SwaggerbondJS contains the Javascript for Swaggerbond
var SwaggerbondJS = `
let searchInput = {
    data: () => ({
            term: ""
        }),
    props: {
        searching: Boolean
    },
    methods: {
        raiseTermEntered: function(e) {
            if (e.keyCode === 13) {
                this.$emit('term-entered', this.term)
            }
        }
    },
    template: ` + "`" + `
        <div class="input-group mb-3">
            <div class="input-group-prepend">              
                <span style="width: 75px" class="input-group-text justify-content-center" id="basic-addon1">
                    <span v-if="searching" class="spinner-grow spinner-grow-sm" role="status"><span class="sr-only">Searching...</span></span>    
                    <span v-if="!searching">Search</span>    
                </span>
            </div>
            <input type="text" v-model="term" @keyup="raiseTermEntered" class="form-control" placeholder="service name, description or * for everything" aria-label="Username" aria-describedby="basic-addon1">
        </div>` + "`" + `
}

let resultsView = {
    props: { 
        results: Object
    },
    template: ` + "`" + `
        <div class="row pt-3 justify-content-left" :class="{ 'justify-content-center': results.data.length === 0 && results.populated }">
            <div v-if="results.data.length === 0 && results.populated" class="font-italic text-muted">Your search did not match any services</div>
            <div class="col-12 col-md-6 col-xl-4 pb-3" v-for="result in results.data">
                <div class="card bg-light" style="height: 300px" :class="{ 'border border-primary': result.highlight && results.mixedRelevance }">
                    <div class="card-body">
                        <h5 class="card-title">{{ result.summary.title }} <small class="text-muted">v{{ result.summary.version }}</small></h5>
                        <p class="card-text overflow-hidden" style="height: 20px">Tags: <span class="font-italic">{{ result.summary.tags.join(', ') }}</span></p>
                        <p class="card-text overflow-hidden" style="height: 160px">{{ result.summary.description }}</p>
                        <a target="_blank" :href="'services/' + result.summary.slug" class="card-link stretched-link">View</a>
                    </div>
                </div>
            </div>
        </div>` + "`" + `
}

let app = new Vue({
    el: "#app",
    data: {
        searching: false,
        results: {
            mixedRelevance: false,
            populated: false,
            data: []
        }
    },
    components: {
        searchInput: searchInput,
        resultsView: resultsView
    },
    methods: {
        searchServices: async function(term) {
            let reset = () => {
                this.searching = false
                this.results.mixedRelevance = false
                this.results.populated = false
                this.results. data = []
            }

            if (term.trim()) {
                reset()
                this.searching = true

                this.results.data = await (await fetch(term === '*' ? 'services/' : ` + "`" + `services?search=${encodeURIComponent(term)}` + "`" + `)).json()

                this.results.data.forEach(r => { 
                    r.highlight = r.relevance == 1 || r.relevance == 2
                    r.summary.tags = r.summary.tags || []
                })
                this.results.mixedRelevance = this.results.data.filter(r => r.highlight).length < this.results.data.length
                this.results.populated = true
                this.searching = false
            } else {
                reset()
            }
        }
    }
})
`
