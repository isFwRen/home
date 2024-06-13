<template>
	<div class="sliding-background">
		<div class="background background1"></div>
		<div class="background background2"></div>
		<div class="background background3"></div>
	</div>
</template>

<style scoped lang="scss">
.sliding-background {
	width: 100%;
	height: 100%;

	.background {
		position: fixed;
		top: 0;
		right: -50%;
		bottom: 0;
		left: -50%;
		opacity: 0.5;
		background-image: linear-gradient(-60deg, #6c3 50%, #09f 50%);
		animation: slide 3s ease-in-out infinite alternate;
	}

	.background2 {
		animation-direction: alternate-reverse;
		animation-duration: 4s;
	}

	.background3 {
		animation-duration: 5s;
	}

	@keyframes slide {
		0% {
			transform: translateX(-25%);
		}
		100% {
			transform: translateX(25%);
		}
	}
}
</style>
