/**
 * Minimal, dependency-free sparkline renderer.
 *
 * Sparkline.draw(svgElement, datapoints, options)
 *   datapoints: array of numbers, or objects like { value, name, date }
 *   options: { stroke, fill, onmousemove(event, datapoint), onmouseout(event) }
 *
 * Stroke/fill fall back to the svg's own stroke/fill attributes, so callers
 * can colour the graph declaratively in markup.
 */
var Sparkline = (function () {
    var NS = 'http://www.w3.org/2000/svg';

    function num(value, fallback) {
        value = parseFloat(value);
        return isNaN(value) ? fallback : value;
    }

    function toValue(point) {
        return (point && typeof point === 'object') ? num(point.value, 0) : num(point, 0);
    }

    function el(name, attrs) {
        var node = document.createElementNS(NS, name);
        for (var key in attrs) {
            if (Object.prototype.hasOwnProperty.call(attrs, key)) {
                node.setAttribute(key, attrs[key]);
            }
        }
        return node;
    }

    function draw(svg, datapoints, options) {
        if (!svg || !datapoints || !datapoints.length) {
            return;
        }
        options = options || {};

        var width = num(svg.getAttribute('width'), 100);
        var height = num(svg.getAttribute('height'), 20);
        var strokeWidth = num(svg.getAttribute('stroke-width'), 2);
        var stroke = options.stroke || svg.getAttribute('stroke') || '#4f46e5';
        var fill = options.fill || svg.getAttribute('fill') || 'none';

        var values = datapoints.map(toValue);
        var min = Math.min.apply(null, values);
        var max = Math.max.apply(null, values);
        var range = (max - min) || 1;
        var pad = strokeWidth + 1;
        var stepX = values.length > 1 ? width / (values.length - 1) : width;

        var points = values.map(function (value, i) {
            return {
                x: i * stepX,
                y: height - pad - ((value - min) / range) * (height - pad * 2)
            };
        });

        var linePath = points.map(function (p, i) {
            return (i === 0 ? 'M' : 'L') + p.x.toFixed(2) + ' ' + p.y.toFixed(2);
        }).join(' ');

        while (svg.firstChild) {
            svg.removeChild(svg.firstChild);
        }

        if (fill && fill !== 'none') {
            var last = points[points.length - 1];
            var first = points[0];
            var area = linePath +
                ' L' + last.x.toFixed(2) + ' ' + height +
                ' L' + first.x.toFixed(2) + ' ' + height + ' Z';
            svg.appendChild(el('path', { d: area, fill: fill, stroke: 'none' }));
        }

        svg.appendChild(el('path', {
            d: linePath,
            fill: 'none',
            stroke: stroke,
            'stroke-width': strokeWidth,
            'stroke-linecap': 'round',
            'stroke-linejoin': 'round'
        }));

        var spotPoint = points[points.length - 1];
        var spot = el('circle', {
            'class': 'sparkline--spot',
            cx: spotPoint.x.toFixed(2),
            cy: spotPoint.y.toFixed(2),
            r: Math.max(1.6, strokeWidth),
            fill: stroke
        });
        svg.appendChild(spot);

        var cursor = el('line', {
            'class': 'sparkline--cursor',
            stroke: stroke,
            x1: 0, x2: 0, y1: 0, y2: height,
            visibility: 'hidden'
        });
        svg.appendChild(cursor);

        if (options.onmousemove || options.onmouseout) {
            svg.addEventListener('mousemove', function (event) {
                var rect = svg.getBoundingClientRect();
                var idx = Math.round((event.clientX - rect.left) / stepX);
                idx = Math.max(0, Math.min(values.length - 1, idx));
                cursor.setAttribute('x1', points[idx].x);
                cursor.setAttribute('x2', points[idx].x);
                cursor.setAttribute('visibility', 'visible');
                if (options.onmousemove) {
                    options.onmousemove(event, datapoints[idx]);
                }
            });
            svg.addEventListener('mouseout', function (event) {
                cursor.setAttribute('visibility', 'hidden');
                if (options.onmouseout) {
                    options.onmouseout(event);
                }
            });
        }
    }

    return { draw: draw };
})();
