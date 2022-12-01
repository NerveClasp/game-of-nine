<script lang="ts">
  import type { CardType } from '../types';

  export let hidden = false;
  export let clickable = false;
  export let card: CardType | undefined = undefined;
</script>

<div
  class={`card ${hidden ? '' : card.kind}`}
  class:hidden
  class:inactive={card.inactive}
  class:playable={card.playable}
  class:clickable
  on:click
>
  {#if !hidden && card}
    <span class="top">{card.value}</span>
    <span class="center">{card.kind}</span>
    <span class="btm">{card.value}</span>
  {/if}
</div>

<style>
  :root {
    --card-bg: antiquewhite;
  }
  .card {
    --card-height: calc(64vh / var(--rows-of-cards, 4));
    --card-width: calc(var(--card-height, 90) * 0.66);
    --abs-padding: calc(var(--card-width) / 8);

    display: flex;
    width: var(--card-width);
    height: var(--card-height);
    background-color: var(--card-bg);
    border: 2px solid var(--card-bg);
    border-radius: 8px;
    font-size: smaller;
    align-items: center;
    justify-content: center;
    margin: 4px;
    position: relative;
    color: #000;
  }

  .top,
  .center,
  .btm {
    position: absolute;
  }

  .top {
    top: var(--abs-padding);
    left: var(--abs-padding);
  }

  .center {
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: xx-large;
  }

  .btm {
    bottom: var(--abs-padding);
    right: var(--abs-padding);
    transform: rotateZ(180deg);
  }

  .inactive {
    opacity: 0.5;
  }

  .playable {
    border-color: red;
  }

  .clickable {
    cursor: pointer;
  }

  .❤,
  .♦ {
    color: red;
  }

  .hidden {
    --card-bg: goldenrod;
    border-color: var(--card-bg);
  }
</style>
