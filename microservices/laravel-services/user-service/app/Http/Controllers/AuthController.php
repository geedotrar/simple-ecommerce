<?php

namespace App\Http\Controllers;

use App\Services\AuthService;
use Illuminate\Support\Facades\Redis;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;
use Tymon\JWTAuth\Exceptions\JWTException;

class AuthController extends Controller
{
    protected $authService;

    public function __construct(AuthService $authService)
    {
        $this->authService = $authService;
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
    
            return response()->json([
                'status' => 'success',
                'message' => 'User registered',
                'data' => $user,
                'error' => null
            ], 201);
        } catch (ValidationException $e) {
            return response()->json([
                'status' => 'error',
                'message' => 'Validation failed',
                'data' => null,
                'error' => $e->errors()
            ], 422);
        } catch (\Exception $e) {
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

                return response()->json([
                    'status' => 'error',
                    'message' => 'Invalid credentials',
                    'data' => null,
                    'error' => 'Email or password is incorrect.'
                ], 401);
            }
            Redis::del($key);

            return response()->json([
                'status' => 'success',
                'message' => 'Login successful',
                'data' => ['token' => $token],
                'error' => null
            ], 200);

        } catch (ValidationException $e) {
            return response()->json([
                'status' => 'error',
                'message' => 'Validation failed',
                'data' => null,
                'error' => $e->errors()
            ], 422);
        } catch (\Exception $e) {
            return response()->json([
                'status' => 'error',
                'message' => 'Login failed',
                'data' => null,
                'error' => $e->getMessage()
            ], 500);
        }
    }

    public function logout(): JsonResponse
    {
        try {
            $this->authService->logout();

            return response()->json([
                'status' => 'success',
                'message' => 'Logout successful',
                'data' => null,
                'error' => null
            ], 200);
        } catch (JWTException $e) {
            return response()->json([
                'status' => 'error',
                'message' => 'Logout failed',
                'data' => null,
                'error' => $e->getMessage()
            ], 500);
        }
    }
}