export class ErrorIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `
            <style>
                .error-icon {
                    stroke: currentColor; /* Use the text color of the parent element for the stroke */
                    flex-shrink: 0; /* Prevent the icon from shrinking in flex layouts */
                    height: 24px; /* Set height to 24 pixels */
                    width: 24px; /* Set width to 24 pixels */
                    fill: white;
                }
            </style>
            <svg xmlns="http://www.w3.org/2000/svg" class="error-icon" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
        `;
   }
}

export class CloseIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `
            <style>
                .close-icon {
                    height: 24px;
                    width: 24px; 
                    fill: none; 
                    stroke: currentColor; /* Use the text color of the parent element for the stroke */
                    stroke-linecap: round;
                    stroke-linejoin: round;
                    stroke-width: 2px;
                }
            </style>
            <svg xmlns="http://www.w3.org/2000/svg" class="close-icon" viewBox="0 0 24 24">
                <path d="M6 18L18 6M6 6l12 12"></path>
            </svg>
        `;
   }
}

export class ExpandIcon extends HTMLElement {
   connectedCallback() {
      this.innerHTML = `
      <svg width="1rem" height="1rem" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M9 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>`;
   }
}

export class CollapseIcon extends HTMLElement {
   connectedCallback() {
      this.innerHTML = `
      <svg width="1rem" height="1rem" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M6 9l6 6 6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>`;
   }
}

export class BlankIcon extends HTMLElement {
   connectedCallback() {
      this.innerHTML = `
      <svg width="1rem" height="1rem" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d=""/>
      </svg>`;
   }
}

customElements.define('error-icon', ErrorIcon);
customElements.define('close-icon', CloseIcon);
customElements.define('expand-icon', ExpandIcon);
customElements.define('collapse-icon', CollapseIcon);
customElements.define('blank-icon', BlankIcon);

