<?php

namespace Database\Seeders;

use App\Models\Role;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Log;

class UserSeeder extends Seeder
{
    public function run(): void
    {
        Log::info('Starting regular user seeding');

        $userRole = Role::firstOrCreate(['name' => 'user']);

        $regularUser = User::firstOrCreate(
            ['email' => 'user@ags.com'],
            [
                'name' => 'User Dummy',
                'password' => Hash::make('password'),
            ]
        );
        $regularUser->roles()->sync([$userRole->id]);

        Log::info('Regular user created and assigned role', ['email' => $regularUser->email, 'role' => $userRole->name]);

        Log::info('Regular user seeding completed');
    }
}
