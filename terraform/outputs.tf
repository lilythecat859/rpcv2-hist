output "instance_ip" {
  description = "Public IP of the instance"
  value       = google_compute_instance.rpcv2_hist.network_interface[0].access_config[0].nat_ip
}