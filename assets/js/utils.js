export class Logger {
  static log(message, type = 'info') {
    const styles = {
      info: 'color: #00ff9d; background: #111; padding: 2px 4px; border-radius: 3px;',
      warn: 'color: #ff9d00; background: #111; padding: 2px 4px; border-radius: 3px;',
      error: 'color: #ff003c; background: #111; padding: 2px 4px; border-radius: 3px;'
    };
    console.log(`%c${message}`, styles[type]);
  }
}

export function debounce(func, wait = 100) {
  let timeout;
  return function(...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), wait);
  };
}

export function throttle(func, limit = 100) {
  let inThrottle;
  return function(...args) {
    if (!inThrottle) {
      func.apply(this, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}