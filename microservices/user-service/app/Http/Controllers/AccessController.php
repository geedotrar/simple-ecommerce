<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Redis;

class AccessController extends Controller
{
    public function getPermissions(): JsonResponse
    {
        $user = auth()->user();
        if (!$user) {
            return response()->json(['error' => 'Unauthorized'], 401);
        }

        $role = $user->roles->pluck('name')->first();

        $cacheKey = "role:{$role}";
        $cached = Redis::get($cacheKey);

        if ($cached) {
            $permissions = json_decode($cached, true);
        } else {
            $permissions = $user->roles->first()->permissions->pluck('name')->toArray();
            Redis::setex($cacheKey, 3600, json_encode($permissions));
        }

        return response()->json([
            'user' => [
                'email' => $user->email,
                'role' => $role,
            ],
            'permissions' => $permissions
        ]);
    }
}
