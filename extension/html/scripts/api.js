
const apiServerUrl = "http://localhost:8642"

export function getApiServerUrl(){
   return apiServerUrl;
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




