import {addRecentlySubmittedLinks} from "./chromeUtils.js";

const apiServerUrl = "http://localhost:8642"

export function getApiServerUrl(){
   return apiServerUrl;
}
function apiLinksWithGroupsEndpoint() {
   return `${getApiServerUrl()}/links/with-groups`
}
function apiLinkByURLEndpoint(url) {
   url = encodeURIComponent(url)
   return `${getApiServerUrl()}/links/by-url/${url}`
}
function apiLinkEndpoint(id) {
   return `${getApiServerUrl()}links/${id}`
}

export function getHttpOptions(method,data) {
   let body,headers
   switch (method) {
      case 'POST':
      case 'PUT':
         body = JSON.stringify(data)
         headers = {
            'Content-Type': 'application/json',
         }
         break;
      case 'GET':
      case 'DELETE':
      default:
         headers={}
   }
   return {
      method: method,
      headers: headers,
      body: body,
   }
}

export async function checkApiHealth(callback) {
   const response = await fetch(`${apiServerUrl}/healthz`).catch( (_)=>{return {ok:false}})
   callback(response.ok)
}

export function apiPutLinkByTab(tab) {
   const link = newLinkFromTab(tab)
   apiPutLink(link)
}

function apiPutLink(link) {
   let endpoint;
   if (link.hasOwnProperty('id') && link.id) {
      endpoint = apiLinkEndpoint(link.id)
   } else  if (link.hasOwnProperty('original_url') && link.original_url) {
      endpoint = apiLinkByURLEndpoint(link.original_url)
   } else {
      throw new Error(`Link has neither '.id' nor '.url' non-empty properties.`)
   }
   fetch(endpoint, getHttpOptions('PUT', link))
         .then(response => {
            response.json()
         })
         .then(data => {
            console.log('Success:', data)
         })
         .catch((error) => {
            console.error('Error:', error)
         });
}

export function newLinkFromTab(tab) {
   return {
      tab_id: tab.id,
      original_url: tab.url,
      title: document.title,
      meta: {
         "description": document.querySelector('meta[name="description"]')?.content,
         "og:title": document.querySelector('meta[property="og:title"]')?.content,
         "og:description": document.querySelector('meta[property="og:description"]')?.content,
         "twitter:title": document.querySelector('meta[name="twitter:title"]')?.content,
         "twitter:description": document.querySelector('meta[name="twitter:description"]')?.content,
         // Add more social and schema.org metadata extraction as required
      }
   };
}

/**
 * Send all current links to the API. To be done periodically.
 */
export function apiPostLinksWithGroups(links) {
   if (links.length === 0) {
      return
   }
   const options = getHttpOptions('POST', links)
   fetch(apiLinksWithGroupsEndpoint(), options)
         .then(async response => {
            console.dir(response)
            if (!response.ok) {
               // Handle HTTP errors (status codes other than 200-299)
               return response.text().then(err => {
                  // Handle any errors (including HTTP errors)
                  let status = `${response.status} - ${response.statusText}`
                  console.log('Error: postLinksWithGroups(): ', status, err);
                  throw new Error(`HTTP Error: ${status}. Details: ${JSON.stringify(err)}`);
               });
            }
            let text = await response.text()
            return text !== "" ? JSON.stringify(text) : {};
         })
         .then(data => {
            addRecentlySubmittedLinks(links)
            console.log('Response from server:', data); // Handle the successful response
         })
         .catch(err => {
            console.log('Error: postLinksWithGroups(): ', err); // Handle any errors (including HTTP errors)
         });
}


