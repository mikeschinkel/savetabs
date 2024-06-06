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
function apiDragDropEndpoint() {
   return `${getApiServerUrl()}/drag-drop`
}
function apiLinkEndpoint(id) {
   return `${getApiServerUrl()}/links/${id}`
}
function apiLabelEndpoint(id) {
   return `${getApiServerUrl()}/links/${id}`
}

export function getHttpOptions(method,data) {
   let body,headers
   switch (method) {
      case 'POST':
         // fallthrough
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

export function newDragDrop(dragType,dragIds,dropType,dropId) {
   for(let index = 0; index < dragIds.length; index++) {
      dragIds[index] = parseInt(dragIds[index],10)
   }
   return {
      drag:{
         type: dragType,
         ids: dragIds
      },
      drop:{
         type: dropType,
         id: Number(dropId)
      }
   }
}

export function apiPostOnDrop(dragDrop) {
   const endpoint = apiDragDropEndpoint();
   const drag = dragDrop.drag.ids.join(',');
   const drop = dragDrop.drop;
   httpPost(endpoint,dragDrop);
   console.log("Drag and Drop:", `${drag.type}:${drag.ids} ==> ${drop.type}:${drop.id}`);
}

function apiPutLink(link) {
   let endpoint;
   if (link.hasOwnProperty('id') && link.id) {
      endpoint = apiLinkEndpoint(link.id)
   } else  if (link.hasOwnProperty('url') && link.url) {
      endpoint = apiLinkByURLEndpoint(link.url)
   } else {
      throw new Error(`Link has neither '.id' nor '.url' non-empty properties.`)
   }
   httpPut(endpoint, link)
}

export function apiPutLabel(label) {
   let endpoint;
   if (label.hasOwnProperty('id') && label.id) {
      endpoint = apiLabelEndpoint(label.id)
      throw new Error(`Label does not have an '.id' property.`)
   }
   httpPut(endpoint, label)
}

function httpPost(endpoint,item) {
   httpRequest('POST',endpoint, item);
}

function httpPut(endpoint,item) {
   httpRequest('PUT',endpoint, item);
}

function httpRequest(method,endpoint,item) {
   fetch(endpoint, getHttpOptions(method, item))
         .then(response => {
            response.json()
         })
         .then(data => {
            console.log('Success:', data)
         })
         .catch((error) => {
            console.log('Error:', error)
         });
}

export function newLinkFromTab(tab) {
   return {
      tab_id: tab.id,
      url: tab.url,
      title: document.title,
      html: tab.html
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


