<?php

namespace App\Services;

use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;

class LogService
{
    public function sendLog(array $payload): void
    {
        try {
            $response = Http::timeout(3)->post(env('USER_LOG_SERVICE_URL'), $payload);

            if ($response->failed()) {
                Log::error('Log Service returned error status', [
                    'status' => $response->status(),
                    'body' => $response->body(),
                    'payload' => $payload,
                ]);
            }
        } catch (\Exception $e) {
            Log::error('Log Service connection error: ' . $e->getMessage(), [
                'payload' => $payload,
            ]);
        }
    }
}
