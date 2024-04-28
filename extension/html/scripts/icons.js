class BaseIcon extends HTMLElement {
   static get observedAttributes() {
      return ['d']; // List attributes to observe for changes
   }

   constructor() {
      super();
      this.attachShadow({ mode: 'open' });
      this.render();
   }

   attributeChangedCallback(name, oldValue, newValue) {
      if (name === 'd' && oldValue !== newValue) {
         this.updatePath(newValue);
      }
   }

   connectedCallback() {
      this.render();
   }

   updatePath(d) {
      const path = this.shadowRoot.querySelector('path');
      if (path) {
         path.setAttribute('d', d);
      }
   }

   render() {
      const d = this.getAttribute('d') || '';
      this.shadowRoot.innerHTML = `
      <style>
          .error-icon {
              stroke: currentColor; /* Use the text color of the parent element for the stroke */
              flex-shrink: 0; /* Prevent the icon from shrinking in flex layouts */
              height: 24px; /* Set height to 24 pixels */
              width: 24px; /* Set width to 24 pixels */
              fill: none;
          }
      </style>
      <svg xmlns="http://www.w3.org/2000/svg" class="error-icon" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${d}"></path>
      </svg>`;
   }
}

export class AlertIcon extends BaseIcon {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `<base-icon d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></base-icon>`;
   }
}
export class ErrorIcon extends BaseIcon {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `<base-icon d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></base-icon>`;
   }
}
export class InfoIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `<base-icon d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></base-icon>`;
   }
}
export class SuccessIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `<base-icon d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></base-icon>`;
   }
}
export class WarningIcon extends HTMLElement {
   constructor() {
      super();
      this.attachShadow({mode: 'open'}).innerHTML = `<base-icon d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></base-icon>`;
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

customElements.define('base-icon', BaseIcon);
customElements.define('alert-icon', AlertIcon);
customElements.define('info-icon', InfoIcon);
customElements.define('success-icon', SuccessIcon);
customElements.define('warning-icon', WarningIcon);
customElements.define('error-icon', ErrorIcon);
customElements.define('close-icon', CloseIcon);
customElements.define('expand-icon', ExpandIcon);
customElements.define('collapse-icon', CollapseIcon);
customElements.define('blank-icon', BlankIcon);

