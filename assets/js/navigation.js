export default class NavigationManager {
    constructor() {
      this.navItems = document.querySelectorAll('nav a');
      this.init();
    }
  
    init() {
      this.setupNavigation();
      this.setupSmoothScrolling();
    }
  
    setupNavigation() {
      this.navItems.forEach(item => {
        item.addEventListener('click', (e) => {
          e.preventDefault();
          this.navItems.forEach(i => i.classList.remove('active'));
          item.classList.add('active');
          
          // Add click effect
          item.style.transform = 'scale(0.95)';
          setTimeout(() => item.style.transform = '', 200);
        });
      });
    }
  
    setupSmoothScrolling() {
      document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
          e.preventDefault();
          const target = document.querySelector(this.getAttribute('href'));
          if (target) {
            target.scrollIntoView({
              behavior: 'smooth',
              block: 'start'
            });
          }
        });
      });
    }
  }