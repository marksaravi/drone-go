package icm20789

func accelerometerFullScale(fsr string) float64 {
	switch fsr {
	case "2g":
		return ACCEL_FULL_SCALE_2G
	case "4g":
		return ACCEL_FULL_SCALE_4G
	case "8g":
		return ACCEL_FULL_SCALE_8G
	case "16g":
		return ACCEL_FULL_SCALE_16G
	default:
		return ACCEL_FULL_SCALE_2G
	}
}
