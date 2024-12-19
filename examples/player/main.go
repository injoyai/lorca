package main

import "github.com/injoyai/lorca"

func main() {
	lorca.Run(&lorca.Config{
		Index: `<video id="my-video" class="video-js" controls preload="auto" width="100%"
poster="https://zhangjikai.com/resource/poster.jpg" data-setup='{"aspectRatio":"16:9"}'>
  <source id="_src" src="https://zhangjikai.com/resource/demo.mp4" type='video/mp4' >
  <p class="vjs-no-js">
  </p>
</video>`,
	})
}
