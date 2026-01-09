package models

const (
	Pending = "pending"
	Success = "success"
	Failed  = "failed"
)

// Deactivate every scan job associated to user-id
func DeactivateScanJobs(userId int) error {
	_, err := DB.Exec("UPDATE scan_jobs SET active=false WHERE user_id=?", userId)
	return err
}

func AddScanJob(userId int, jobId string) error {
	_, err := DB.Exec("INSERT INTO scan_jobs(user_id, job_id) VALUES(?, ?)", userId, jobId)
	return err
}

func GetActiveScanJob(userId int) (string, error) {
	var jobId string

	row := DB.QueryRow("SELECT job_id FROM scan_jobs WHERE user_id=? ORDER BY created_at DESC", userId)
	if err := row.Scan(&jobId); err != nil {
		return "", err
	}

	return jobId, nil
}

func ListScanJobs(userId int) ([]string, error) {
	r, err := DB.Query("SELECT job_id FROM scan_jobs WHERE user_id=?", userId)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for r.Next() {
		var jobId string
		r.Scan(&jobId)

		jobIds = append(jobIds, jobId)
	}

	return jobIds, nil
}

func GetScanCountOfUser(userId int) (int, error) {
	var cnt int

	err := DB.QueryRow("SELECT COUNT(*) FROM scan_jobs WHERE user_id=?", userId).Scan(&cnt)
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}
