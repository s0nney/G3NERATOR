import CircuitAnimator from './animations.js';

document.addEventListener('DOMContentLoaded', () => {
    // Initialize circuit animations
    new CircuitAnimator();

    // Header scroll effect
    const header = document.querySelector('header');
    window.addEventListener('scroll', () => {
        header.style.boxShadow = window.scrollY > 50
            ? '0 5px 15px rgba(0, 0, 0, 0.3)'
            : 'none';
    });

    // Article hover effects
    document.querySelectorAll('article').forEach(article => {
        article.addEventListener('mouseenter', () => {
            article.style.transform = 'translateY(-5px)';
            article.style.boxShadow = '0 5px 15px rgba(255, 255, 255, 0.1)';
            article.style.transition = 'all 0.3s ease';
        });

        article.addEventListener('mouseleave', () => {
            article.style.transform = '';
            article.style.boxShadow = '';
        });
    });
});

// Scroll progress bar
window.onscroll = function () {
    var winScroll = document.body.scrollTop || document.documentElement.scrollTop;
    var height = document.documentElement.scrollHeight - document.documentElement.clientHeight;
    var scrolled = (winScroll / height) * 100;
    document.getElementById("progressBar").style.width = scrolled + "%";
};

// Generate TOC from headings
document.addEventListener('DOMContentLoaded', function () {
    const tocList = document.getElementById('toc-list');
    const headings = document.querySelectorAll('.article-content h2, .article-content h3');

    headings.forEach(heading => {
        const id = heading.textContent.toLowerCase().replace(/[^\w]+/g, '-');
        heading.id = id;

        const li = document.createElement('li');
        const a = document.createElement('a');
        a.href = '#' + id;
        a.textContent = heading.textContent;

        if (heading.tagName === 'H3') {
            const subUl = document.createElement('ul');
            const subLi = document.createElement('li');
            subLi.appendChild(a);
            subUl.appendChild(subLi);
            li.appendChild(subUl);
        } else {
            li.appendChild(a);
        }

        tocList.appendChild(li);
    });
});