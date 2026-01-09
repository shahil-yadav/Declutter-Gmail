package models

func AddTrashJob(jobId string, userId int) error {
	_, err := DB.Exec("INSERT INTO trash_jobs(job_id, user_id) VALUES(?, ?)", jobId, userId)
	return err
}

func GetActiveTrashJob(userId int) (string, error) {
	var jobId string

	row := DB.QueryRow("SELECT job_id FROM trash_jobs WHERE user_id=? ORDER BY created_at DESC", userId)
	if err := row.Scan(&jobId); err != nil {
		return "", err
	}

	return jobId, nil
}

func ListTrashJobs(userId int) ([]string, error) {
	var jobs []string
	rows, err := DB.Query("SELECT job_id FROM trash_jobs")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var job string
		rows.Scan(&job)

		jobs = append(jobs, job)
	}

	return jobs, nil
}
