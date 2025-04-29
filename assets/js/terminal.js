export default class TerminalEffect {
  constructor() {
    this.terminalElements = document.querySelectorAll('.terminal, code');
    this.init();
  }

  init() {
    if (this.terminalElements.length > 0) {
      this.setupTerminalEffects();
      this.typeWriterEffect();
    }
  }

  setupTerminalEffects() {
    this.terminalElements.forEach(terminal => {
      // Add blinking cursor
      const cursor = document.createElement('span');
      cursor.className = 'terminal-cursor';
      cursor.textContent = '_';
      terminal.appendChild(cursor);
      
      // Blink animation
      setInterval(() => {
        cursor.style.visibility = cursor.style.visibility === 'hidden' ? 'visible' : 'hidden';
      }, 500);
    });
  }

  typeWriterEffect() {
    const typeElements = document.querySelectorAll('.typewriter');
    typeElements.forEach(el => {
      const text = el.textContent;
      el.textContent = '';
      
      let i = 0;
      const speed = Math.random() * 50 + 50;
      const typing = setInterval(() => {
        if (i < text.length) {
          el.textContent += text.charAt(i);
          i++;
        } else {
          clearInterval(typing);
        }
      }, speed);
    });
  }
}