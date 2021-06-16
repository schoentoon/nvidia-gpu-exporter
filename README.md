# NVidia GPU Exporter

This is yet another nvidia gpu metrics exporter for prometheus.
Main difference with most other exporters for this same purpose is that is uses the actual prometheus golang library rather than building their our own metrics page.
Besides that this library is set up to be more extendable by default, as it's fairly easy to add new metrics to it.
To add new metrics I would like to refer to [collector.go](./collector.go), documentation in the comments should be clear enough.

## Build

To build this project you'll need to have [golang](https://golang.org/) installed.
After that just clone this project and build it with `go build .` while in the directory of your local clone.

## Usage

You can either run the provided binary as is, in which case it will listen on `127.0.0.1:9101` by default.
Or you can choose to override this using the `-address` argument.
For automatically starting upon system boot a systemd service file is provided.
You can either copy this to `~/.config/systemd/user` (I would recommend this in case of desktop usage) or `/etc/systemd/system/`.
Afterwards just run `systemctl daemon-reload` (add `--user` in case of the user directory).
And `systemctl enable --now nvidia-gpu-exporter` (again add `--user` in case of the user directory), to enable it and start it.
For non-systemd systems you will have to unfortunately write your own startup script or alike, but most likely if you're running a non-systemd system in 2021 you'll know how to do so.

## Prometheus

Finally you'll want to configure prometheus to actually query the exporter.
This is basically just adding the following to your prometheus config (likely located at `/etc/prometheus/prometheus.yml`).
Of course in case you decide to run on a different port you'll have to adjust this.

```yml
scrape_configs:
  - job_name: "gpu_exporter"
    static_configs:
      - targets: ['127.0.0.1:9101']
```

## Exported metrics

And finally, here's a brief example of what metrics you can expect.

```bash
# HELP gpu_driver The version of the installed NVIDIA display driver. This is an alphanumeric string.
# TYPE gpu_driver gauge
gpu_driver{driver="460.80"} 1
# HELP gpu_fan_speed The fan speed value is the percent of the product's maximum noise tolerance fan speed that the device's fan is currently intended to run at.
# TYPE gpu_fan_speed gauge
gpu_fan_speed{name="GeForce GTX 750",uuid="GPU-<snip>"} 30
# HELP gpu_graphics_clock_speed Current frequency of graphics (shader) clock. In megahertz
# TYPE gpu_graphics_clock_speed gauge
gpu_graphics_clock_speed{name="GeForce GT 1030",uuid="GPU-<snip>"} 1721
gpu_graphics_clock_speed{name="GeForce GTX 750",uuid="GPU-<snip>"} 135
# HELP gpu_memory_clock_speed Current frequency of memory clock. In megahertz
# TYPE gpu_memory_clock_speed gauge
gpu_memory_clock_speed{name="GeForce GT 1030",uuid="GPU-<snip>"} 1050
gpu_memory_clock_speed{name="GeForce GTX 750",uuid="GPU-<snip>"} 405
# HELP gpu_memory_free Total free memory. In MiB.
# TYPE gpu_memory_free gauge
gpu_memory_free{name="GeForce GT 1030",uuid="GPU-<snip>"} 467
gpu_memory_free{name="GeForce GTX 750",uuid="GPU-<snip>"} 974
# HELP gpu_memory_total Total installed GPU memory. In MiB.
# TYPE gpu_memory_total gauge
gpu_memory_total{name="GeForce GT 1030",uuid="GPU-<snip>"} 1998
gpu_memory_total{name="GeForce GTX 750",uuid="GPU-<snip>"} 981
# HELP gpu_memory_used Total memory allocated by active contexts. In MiB.
# TYPE gpu_memory_used gauge
gpu_memory_used{name="GeForce GT 1030",uuid="GPU-<snip>"} 1531
gpu_memory_used{name="GeForce GTX 750",uuid="GPU-<snip>"} 7
# HELP gpu_power_draw The last measured power draw for the entire board, in watts. Only available if power management is supported. This reading is accurate to within +/- 5 watts.
# TYPE gpu_power_draw gauge
gpu_power_draw{name="GeForce GTX 750",uuid="GPU-<snip>"} 0.63
# HELP gpu_sm_clock_speed Current frequency of SM (Streaming Multiprocessor) clock. In megahertz
# TYPE gpu_sm_clock_speed gauge
gpu_sm_clock_speed{name="GeForce GT 1030",uuid="GPU-<snip>"} 1721
gpu_sm_clock_speed{name="GeForce GTX 750",uuid="GPU-<snip>"} 135
# HELP gpu_temperature Core GPU temperature. in degrees C.
# TYPE gpu_temperature gauge
gpu_temperature{name="GeForce GT 1030",uuid="GPU-<snip>"} 44
gpu_temperature{name="GeForce GTX 750",uuid="GPU-<snip>"} 30
# HELP gpu_utilization Percent of time over the past sample period during which one or more kernels was executing on the GPU.
# TYPE gpu_utilization gauge
gpu_utilization{name="GeForce GT 1030",uuid="GPU-<snip>"} 78
gpu_utilization{name="GeForce GTX 750",uuid="GPU-<snip>"} 0
# HELP gpu_video_clock_speed Current frequency of video encoder/decoder clock. In megahertz
# TYPE gpu_video_clock_speed gauge
gpu_video_clock_speed{name="GeForce GT 1030",uuid="GPU-<snip>"} 1544
gpu_video_clock_speed{name="GeForce GTX 750",uuid="GPU-<snip>"} 405
```
