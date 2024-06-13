<template>
  <div class="z-flex justify-center spinners">
    <flower-spinner
      v-if="randomNum === 1"
      :animation-duration="2500"
      :size="70"
      :color="color"
    />

    <pixel-spinner
      v-else-if="randomNum === 2"
      :animation-duration="2000"
      :pixel-size="70"
      :color="color"
    />

    <hollow-dots-spinner
      v-else-if="randomNum === 3"
      :animation-duration="1000"
      :dot-size="15"
      :dots-num="3"
      :color="color"
    />

    <intersecting-circles-spinner
      v-else-if="randomNum === 4"
      :animation-duration="1200"
      :size="70"
      :color="color"
    />

    <orbit-spinner
      v-else-if="randomNum === 5"
      :animation-duration="1200"
      :size="55"
      :color="color"
    />

    <radar-spinner
      v-else-if="randomNum === 6"
      :animation-duration="2000"
      :size="60"
      :color="color"
    />

    <scaling-squares-spinner
      v-else-if="randomNum === 7"
      :animation-duration="1250"
      :size="65"
      :color="color"
    />

    <half-circle-spinner
      v-else-if="randomNum === 8"
      :animation-duration="1000"
      :size="60"
      :color="color"
    />

    <trinity-rings-spinner
      v-else-if="randomNum === 9"
      :animation-duration="1500"
      :size="66"
      :color="color"
    />

    <fulfilling-square-spinner
      v-else-if="randomNum === 10"
      :animation-duration="4000"
      :size="50"
      :color="color"
    />

    <circles-to-rhombuses-spinner
      v-else-if="randomNum === 11"
      :animation-duration="1200"
      :circles-num="3"
      :circle-size="15"
      :color="color"
    />
        
    <semipolar-spinner
      v-else-if="randomNum === 12"
      :animation-duration="2000"
      :size="65"
      :color="color"
    />
        
    <self-building-square-spinner
      v-else-if="randomNum === 13"
      :animation-duration="6000"
      :size="40"
      :color="color"
    />
        
    <swapping-squares-spinner
      v-else-if="randomNum === 14"
      :animation-duration="1000"
      :size="65"
      :color="color"
    />

    <fulfilling-bouncing-circle-spinner
      v-else-if="randomNum === 15"
      :animation-duration="4000"
      :size="60"
      :color="color"
    />

    <fingerprint-spinner
      v-else-if="randomNum === 16"
      :animation-duration="1500"
      :size="64"
      :color="color"
    />

    <spring-spinner
      v-else-if="randomNum === 17"
      :animation-duration="3000"
      :size="60"
      :color="color"
    />

    <atom-spinner
      v-else-if="randomNum === 18"
      :animation-duration="1000"
      :size="60"
      :color="color"
    />

    <looping-rhombuses-spinner
      v-else-if="randomNum === 19"
      :animation-duration="2500"
      :rhombus-size="15"
      :color="color"
    />

    <breeding-rhombus-spinner
      v-else-if="randomNum === 20"
      :animation-duration="2000"
      :size="65"
      :color="color"
    />
  </div>
</template>

<script>
  import { 
    HollowDotsSpinner, 
    PixelSpinner, 
    FlowerSpinner, 
    IntersectingCirclesSpinner, 
    OrbitSpinner, 
    FingerprintSpinner, 
    TrinityRingsSpinner, 
    FulfillingSquareSpinner, 
    CirclesToRhombusesSpinner, 
    SemipolarSpinner,
    BreedingRhombusSpinner, 
    SwappingSquaresSpinner, 
    ScalingSquaresSpinner, 
    FulfillingBouncingCircleSpinner, 
    RadarSpinner,
    SelfBuildingSquareSpinner, 
    SpringSpinner, 
    LoopingRhombusesSpinner, 
    HalfCircleSpinner, 
    AtomSpinner
  } from 'epic-spinners'

  const randomNum = Math.ceil(Math.random()*20)

  export default {
    data() {
      return {
        randomNum,
        color: '#fff'
      }
    },

    components: {
      HollowDotsSpinner, 
      PixelSpinner, 
      FlowerSpinner, 
      IntersectingCirclesSpinner, 
      OrbitSpinner, 
      FingerprintSpinner, 
      TrinityRingsSpinner, 
      FulfillingSquareSpinner, 
      CirclesToRhombusesSpinner, 
      SemipolarSpinner,
      BreedingRhombusSpinner, 
      SwappingSquaresSpinner, 
      ScalingSquaresSpinner, 
      FulfillingBouncingCircleSpinner, 
      RadarSpinner,
      SelfBuildingSquareSpinner, 
      SpringSpinner, 
      LoopingRhombusesSpinner, 
      HalfCircleSpinner, 
      AtomSpinner
    }
  }
</script>