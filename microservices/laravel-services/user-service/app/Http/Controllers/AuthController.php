<?php

namespace App\Http\Controllers;

use App\Services\AuthService;
use App\Services\LogService;
use Illuminate\Support\Facades\Redis;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;
use Tymon\JWTAuth\Exceptions\JWTException;

class AuthController extends Controller
{
    protected $authService;
    protected $logService;

    public function __construct(AuthService $authService, LogService $logService)
    {
        $this->authService = $authService;
        $this->logService = $logService;
    }

    public function register(Request $request): JsonResponse
    {
        try {
            $validated = $request->validate([
                'name' => 'required|string|max:255',
                'email' => 'required|string|email|unique:users',
                'password' => [
                    'required',
                    'string',
                    'min:8',
                    'confirmed',
                    'regex:/^(?=.*[0-9])(?=.*[\W_]).+$/'
                ],
            ], [
                'name.required' => 'The name field is required.',
                'email.required' => 'The email field is required.',
                'email.email' => 'The email format is invalid.',
                'email.unique' => 'The email is already registered.',
                'password.required' => 'The password field is required.',
                'password.min' => 'The password must be at least 8 characters.',
                'password.confirmed' => 'The password confirmation does not match.',
                'password.regex' => 'The password must contain at least have number and have symbol.'
            ]);
    
            $validated['role'] = 'user';
    
            $user = $this->authService->register($validated);

            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_register',
                'user' => [
                    'id' => (string) $user->id,
                    'email' => auth()->user()->email, 
                ],
                'ip_address' => $request->ip(),
                'session_id' => session()->getId(),
                'details' => [
                    'message' => 'User registered',
                    'email' => $user->email,
                    'user_agent' => $request->userAgent(),
                ]
            ]);
    
            return response()->json([
                'status' => 'success',
                'message' => 'User registered',
                'data' => $user,
                'error' => null
            ], 201);
        } catch (ValidationException $e) {
            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_register_failed',
                'user' => null,
                'ip_address' => $request->ip(),
                'session_id' => session()->getId(),
                'details' => [
                    'message' => 'Validation failed',
                    'errors' => $e->errors(),
                    'user_agent' => $request->userAgent(),
                ]
            ]);
            
            return response()->json([
                'status' => 'error',
                'message' => 'Validation failed',
                'data' => null,
                'error' => $e->errors()
            ], 422);
        } catch (\Exception $e) {
            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_register_failed',
                'user' => null,
                'ip_address' => $request->ip(),
                'session_id' => session()->getId(),
                'details' => [
                    'message' => 'Registration failed: ' . $e->getMessage(),
                    'user_agent' => $request->userAgent(),
                ]
            ]);

            return response()->json([
                'status' => 'error',
                'message' => 'Registration failed',
                'data' => null,
                'error' => $e->getMessage()
            ], 500);
        }
    }
    
    public function login(Request $request): JsonResponse
    {
        $ip = $request->ip();
        $key = "login_attempts:{$ip}";
        $maxAttempts = 5;
        $decaySeconds = 60;

        $attempts = Redis::get($key);
        $attempts = $attempts ? (int)$attempts : 0;

        if ($attempts >= $maxAttempts) {
            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_login_failed',
                'user' => null,
                'ip_address' => $ip,
                'session_id' => session()->getId(),
                'details' => [
                    'login_status' => 'failed',
                    'message' => 'Too many login attempts. Please try again later.',
                    'user_agent' => $request->userAgent()
                ]
            ]);

            return response()->json([
                'status' => 'error',
                'message' => 'Too many login attempts. Please try again later.',
                'data' => null,
                'error' => 'Rate limit exceeded.'
            ], 429);
        }

        try {
            $credentials = $request->validate([
                'email' => 'required|email',
                'password' => 'required'
            ]);

            $token = $this->authService->login($credentials);

            if (!$token) {
                if ($attempts === 0) {
                    Redis::setex($key, $decaySeconds, 1);
                } else {
                    Redis::incr($key);
                }

                $this->logService->sendLog([
                    'timestamp' => now()->toIso8601String(),
                    'event' => 'user_login_failed',
                    'user' => null,
                    'ip_address' => $ip,
                    'session_id' => session()->getId(),
                    'details' => [
                        'login_status' => 'failed',
                        'message' => 'Invalid credentials',
                        'user_agent' => $request->userAgent()
                    ]
                ]);


                return response()->json([
                    'status' => 'error',
                    'message' => 'Invalid credentials',
                    'data' => null,
                    'error' => 'Email or password is incorrect.'
                ], 401);
            }
            Redis::del($key);

            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_login',
                'user' => [
                    'id' => (string) auth()->user()->id,
                    'email' => auth()->user()->email, 
                ],
                'ip_address' => $request->ip(),
                'session_id' => session()->getId(), 
                'details' => [
                    'login_status' => 'success',
                    'message' => 'Login successful',
                    'user_agent' => $request->userAgent()
                ]
            ]);

            return response()->json([
                'status' => 'success',
                'message' => 'Login successful',
                'data' => ['token' => $token],
                'error' => null
            ], 200);

        } catch (ValidationException $e) {
            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_login_failed',
                'user' => null,
                'ip_address' => $ip,
                'session_id' => session()->getId(),
                'details' => [
                    'login_status' => 'failed',
                    'message' => 'Validation failed',
                    'errors' => $e->errors(),
                    'user_agent' => $request->userAgent()
                ]
            ]);

            return response()->json([
                'status' => 'error',
                'message' => 'Validation failed',
                'data' => null,
                'error' => $e->errors()
            ], 422);
        } catch (\Exception $e) {
            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_login_failed',
                'user' => null,
                'ip_address' => $ip,
                'session_id' => session()->getId(),
                'details' => [
                    'login_status' => 'failed',
                    'message' => 'Login failed: ' . $e->getMessage(),
                    'user_agent' => $request->userAgent()
                ]
            ]);

            return response()->json([
                'status' => 'error',
                'message' => 'Login failed',
                'data' => null,
                'error' => $e->getMessage()
            ], 500);
        }
    }

    public function logout(Request $request): JsonResponse
    {
        try {
            $this->authService->logout();

            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_logout',
                'user' => [
                    'id' => (string) auth()->user()->id,
                    'email' => auth()->user()->email, 
                ],
                'ip_address' => $request->ip(),
                'session_id' => session()->getId(),
                'details' => [
                    'logout_status' => 'success',
                    'message' => 'Logout successful',
                    'user_agent' => $request->userAgent(),
                ]
            ]);

            return response()->json([
                'status' => 'success',
                'message' => 'Logout successful',
                'data' => null,
                'error' => null
            ], 200);
        } catch (JWTException $e) {
            $this->logService->sendLog([
                'timestamp' => now()->toIso8601String(),
                'event' => 'user_logout_failed',
                'user' => auth()->user() ? [
                    'id' => (string) auth()->user()->id,
                    'email' => auth()->user()->email,
                ] : null,
                'ip_address' => $request->ip(),
                'session_id' => session()->getId(),
                'details' => [
                    'logout_status' => 'failed',
                    'message' => 'Logout failed',
                    'error' => $e->getMessage(),
                    'user_agent' => $request->userAgent(),
                ]
            ]);


            return response()->json([
                'status' => 'error',
                'message' => 'Logout failed',
                'data' => null,
                'error' => $e->getMessage()
            ], 500);
        }
    }
}