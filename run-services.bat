@echo off
echo Starting Mopcare LMS Microservices...

start "Course Service" cmd /k "cd /d %~dp0services\course-service && go run main.go"
timeout /t 2 /nobreak >nul

start "User Service" cmd /k "cd /d %~dp0services\user-service && go run main.go"
timeout /t 2 /nobreak >nul

start "Enrollment Service" cmd /k "cd /d %~dp0services\enrollment-service && go run main.go"
timeout /t 2 /nobreak >nul

start "API Gateway" cmd /k "cd /d %~dp0gateway-fiber && go run main.go"

echo All services started!
echo API Gateway available at: http://localhost:9090
pause