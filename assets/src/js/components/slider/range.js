import * as $ from "bm.js/bm.module.js";
import {html, render} from 'lit-html';

var tmpl = (elem) => html`
	<style>
		@import "/static/css/root.css";

		:host {
			width: ${elem.attr("width")||"7rem"};
		}

		[flex] {
			display: flex;
			justify-content: space-between;
		}

		#track {
			height: 1rem;

			position: relative;

			input[type=range] {
				position: absolute;
				top:    0;
				bottom: 0;
				width: 100%;
				margin: 0;

				appearance: none;
				background-color: transparent;
				pointer-events: none;

				&::-webkit-slider-thumb {
					appearance: none;
					height: 1rem;
					width: 0.5rem;
					pointer-events: auto;
					background: var(--blue-400);
				}
			}
			& #bar {
				height: 5px;
				width: 100%;
				background: gray;

				position: absolute;
				top: 0;
				bottom: 0;
				margin: auto;

				& > div {
					position: absolute;
					top: 0;
					bottom: 0;

					background: var(--blue-400);
				}
			}
			input[type=range][role=high]::-webkit-slider-thumb {
				border-radius: 0 50% 50% 0;
				transform: translateX(50%);
				border-left: 1px solid gray;
			}
			input[type=range][role=low]::-webkit-slider-thumb {
				border-radius: 50% 0 0 50%;
				transform: translateX(-50%);
				border-right: 1px solid gray;
			}
		}
		input[type=number] {
			width: 3rem;
		}
 	</style>
	<section flex>
		<span>${elem.min}</span>
		<span>${elem.max}</span>
	</section>
	<section id="track">
		<div id="bar"><div></div></div>
		<input type="range" role="low"  min="${elem.min}" max="${elem.max}" value="${elem.attr("start") || elem.min}" step="${elem.step}" @input="${evt => elem.onRangeChange(evt)}" />
		<input type="range" role="high" min="${elem.min}" max="${elem.max}" value="${elem.attr("end")   || elem.max}" step="${elem.step}" @input="${evt => elem.onRangeChange(evt)}" />
	</section>
	<section>
			<input type="number" role="low"  min="${elem.min}" max="${elem.max}" value="${elem.attr("start") || elem.min}" step="${elem.step}" @input="${evt => elem.onInputChange(evt)}" />
			<span> - </span>
			<input type="number" role="high" min="${elem.min}" max="${elem.max}" value="${elem.attr("end")   || elem.max}" step="${elem.step}" @input="${evt => elem.onInputChange(evt)}" value="5"/>
	</section>
`;

class RangeSlider extends $.CustomElement {
	constructor() {
		super();
	}

	static get observedAttributes() {
		return [];
	}
	onAttributeChanged(name, oValue, nValue) {
	}

	get start() {
		return $.get(this.shadowRoot, "input[role=low]").value;
	}
	get end() {
		return $.get(this.shadowRoot, "input[role=high]").value;
	}
	get max() {
		return Number(this.attr("max")) || 100;
	}
	get min() {
		return Number(this.attr("min")) || 0;
	}
	get step() {
		return Number(this.attr("step")) || 1;
	}

	onConnected(){
		this.updateBar(this.start, this.end, this.min, this.max);
	}

	async render() {
		render(tmpl(this), this.shadowRoot);
	}

	//handler
	onRangeChange(evt) {
		$.get(this.shadowRoot, `input[type=number][role=${evt.target.attr("role")}]`).value = evt.target.value;

		this.updateBar(this.start, this.end, this.min, this.max);
	}
	onInputChange(evt) {
		$.get(this.shadowRoot, `input[type=range][role=${evt.target.attr("role")}]`).value = evt.target.value;
		this.updateBar(this.start, this.end, this.min, this.max);
	}
	updateBar(start, end, min, max) {
		let total = max - min;
		let $color = $.get(this.shadowRoot, "#bar div");

		$color.style.left = `${((start - min)/total)*100}%`;
		$color.style.right = `${((max - end)/total)*100}%`;
		console.log(total, start, end, min, max,(start - min)/total,(max - end)/total);
	}
}
customElements.define("c-range-slider", RangeSlider);
