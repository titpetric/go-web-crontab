var Sparkline = (function() {
    function num(v, d) {
        v = parseFloat(v);
        if (isNaN(v)) { return d; }
        return v;
    }
    function draw(svg, values, options) {
        if (!svg || !values || values.length === 0) { return; }
        options = options || {};
        var width = num(svg.getAttribute('width'), 100);
        var height = num(svg.getAttribute('height'), 20);
        var stroke = options.stroke || 'blue';
        var strokeWidth = num(svg.getAttribute('stroke-width'), 2);
        var min = values[0];
        var max = values[0];
        var i;
        for (i = 0; i < values.length; i++) {
            values[i] = num(values[i], 0);
            if (values[i] < min) { min = values[i]; }
            if (values[i] > max) { max = values[i]; }
        }
        var range = max - min;
        if (range === 0) { range = 1; }
        var step = width;
        if (values.length > 1) { step = width / (values.length - 1); }
        var d = '';
        for (i = 0; i < values.length; i++) {
            var x = i * step;
            var y = height - ((values[i] - min) / range * (height - strokeWidth * 2)) - strokeWidth;
            d += (i === 0 ? 'M' : 'L') + x.toFixed(2) + ' ' + y.toFixed(2) + ' ';
        }
        while (svg.firstChild) { svg.removeChild(svg.firstChild); }
        var path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        path.setAttribute('d', d);
        path.setAttribute('fill', 'none');
        path.setAttribute('stroke', stroke);
        path.setAttribute('stroke-width', strokeWidth);
        path.setAttribute('stroke-linecap', 'round');
        path.setAttribute('stroke-linejoin', 'round');
        svg.appendChild(path);
    }
    return { draw: draw };
})();
