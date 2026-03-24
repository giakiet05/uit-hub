
# API Endpoints

## GET

### Student

- `/student/profile`
- `/student/schedule?year=&semester=`
- `/student/schedule/exam?year=&semester=`
- `/student/score`
- `/student/lookup/tuitionfee`
- `/student/insurance`
- `/student/courses`
- `/student/training-points`
- `/student/lookup/office365`
- `/student/survey-form`
- `/student/deadlines`
- `/student/courses/{courseId}/materials`
- `/student/courses/{courseId}/assignments`

### Other

- `/annual-plan`
- `/program/{slugs}`
- `/tutorial`

### Regulations

- `/regulations/{slugs}`
- `/regulations-vnu/{slugs}`
- `/regulations-moet/{slugs}`
- `/regulations-short-term`
- `/regulations/procedures/lecturer`
- `/regulations/procedures/student`

### Notifications

- `/notifications`

### Rooms

- `/rooms/availability?date=&start=&end=`

## POST

### Auth

- `/login`

### Student

- `/student/transcript-regis`
- `/student/referral`
- `/student/tuition-extend`
- `/student/outpatient-regis`
- `/student/eor-regis`
- `/student/monthly-parking`
- `/student/graduate`
- `/student/graduation-thesis`
- `/student/assignments/{assignmentId}/submissions`

### Other

- `/contact`
- `/verification`