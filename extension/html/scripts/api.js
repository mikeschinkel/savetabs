
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