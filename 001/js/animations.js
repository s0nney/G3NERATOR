class CircuitAnimator {
    constructor() {
        this.circuitElements = [];
        this.init();
    }

    init() {
        this.createInitialCircuits();
        setInterval(() => this.createCircuitElement(), 3000);
    }

    createCircuitElement() {
        const circuit = document.createElement('div');
        circuit.className = 'circuit-line-animated';
        
        const size = Math.random() * 2 + 1;
        const isHorizontal = Math.random() > 0.5;
        
        Object.assign(circuit.style, {
            position: 'fixed',
            backgroundColor: 'rgba(255, 255, 255, 0.07)',
            zIndex: '-1',
            opacity: '0',
            transition: 'opacity 2s'
        });

        if (isHorizontal) {
            circuit.style.height = `${size}px`;
            circuit.style.width = `${Math.random() * 100 + 50}px`;
        } else {
            circuit.style.width = `${size}px`;
            circuit.style.height = `${Math.random() * 100 + 50}px`;
        }

        circuit.style.top = `${Math.random() * 100}vh`;
        circuit.style.left = `${Math.random() * 100}vw`;

        document.body.appendChild(circuit);
        this.circuitElements.push(circuit);

        setTimeout(() => {
            circuit.style.opacity = '1';
            setTimeout(() => {
                circuit.style.opacity = '0';
                setTimeout(() => {
                    circuit.remove();
                    this.circuitElements = this.circuitElements.filter(c => c !== circuit);
                }, 2000);
            }, 4000);
        }, 100);
    }

    createInitialCircuits() {
        for (let i = 0; i < 5; i++) {
            setTimeout(() => this.createCircuitElement(), i * 1000);
        }
    }
}

export default CircuitAnimator;