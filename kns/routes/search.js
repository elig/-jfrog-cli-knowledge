/* The api endpoints to fetch the data from algolia */
require('dotenv').config();

const express = require('express'),
    router = express.Router(),
    algoliasearch = require('algoliasearch');


// Define the algolia credentials
const client = algoliasearch(process.env.ALGOLIA_APPID, process.env.ALGOLIA_APIKEY),
    index = client.initIndex('wp_searchable_posts');



/*
**  Ping know node service
**
**  GET /api/v1/search/ping
*/
router.get('/ping', (req, res) => {
    res.send({ 'message': 'pong' });
});

/*
**  Get the facets list
**
**  GET /api/v1/search/facets?query={searchQuery}
*/
router.get('/facets', (req, res) => {
    index.search(req.query.query, {
        facets: ['*'],
    }).then((response) => {

        if (typeof response.facets !== 'undefined') {
            // get only the post_type_label
            let facetsOutput = response.facets.post_type_label;

            if (typeof facetsOutput !== 'undefined') {
                // remove the mega menu from the facets
                delete facetsOutput['Mega Menus'];
                res.send(facetsOutput);
            } else {
                res.send([]);
            }
        } else {
            res.send([]);
        }

    }).catch(err => {
        console.error(err);
        res.send([]);
    });
});

/*
**  Get the actual results from Algolia with the search query string and the facet for the post type
**  GET /api/v1/search?query={queryString}&facet={facetName}
**  @params - query, facet
*/
router.get('/', (req, res) => {
    // make the request to algolia
    index.search(req.query.query, {
        facets: ['post_type_label'],
        facetFilters: [
            [`post_type_label:${req.query.facet}`],
            ['locale:en_US']
        ],
        hitsPerPage: 10,
        attributesToRetrieve: [
            'post_id',
            'permalink',
            'post_title',
            'post_author',
            'post_date',
            'post_date_formatted',
            'post_type_label',
        ]
    }).then(({ hits }) => {

        let output = [];

        if (hits.length) {
            hits.forEach(element => {
                // return null for author if not exists
                let postAuthor = (typeof element.post_author !== 'undefined') ? element.post_author.display_name : null;

                // add the data to the output
                output.push({
                    post_id: element.post_id,
                    url: element.permalink,
                    title: element.post_title,
                    author: postAuthor,
                    publish_date: element.post_date_formatted,
                    content_type: element.post_type_label,
                });
            });

            res.send(output);
        } else {
            res.statusCode = 404;
            res.send([]);
        }
    }).catch(err => {
        console.error(err);
        res.statusCode = 500;
        res.send(err);
    });
});


/*
**  Get the individual post by id
**  GET /api/v1/search/id/{postId}
**  @params - query, facet, postId
*/
router.get('/id/:postId', (req, res) => {
    // make the request to algolia
    index.search('', {
        hitsPerPage: 10,
        numericFilters: [
            `post_id=${req.params.postId}`
        ],
        attributesToRetrieve: [
            'post_id',
            'permalink',
            'post_title',
            'post_author',
            'post_date',
            'post_date_formatted',
            'post_type_label',
            'content'
        ]
    }).then(({ hits }) => {
        if (hits.length === 1) {
            // return null for author if not exists
            let postAuthor = (typeof hits[0].post_author !== 'undefined') ? hits[0].post_author.display_name : null;

            let output = {
                post_id: hits[0].post_id,
                url: hits[0].permalink,
                title: hits[0].post_title,
                author: postAuthor,
                publish_date: hits[0].post_date_formatted,
                content_type: hits[0].post_type_label,
                content: hits[0].content
            };

            res.send(output);
        } else {
            res.statusCode = 404;
            res.send([]);
        }
    }).catch(err => {
        console.error(err);
        res.statusCode = 500;
        res.send([err]);
    });
});

module.exports = router;