<?php

namespace Database\Seeders;

use App\Models\Role;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Log;

class AdminSeeder extends Seeder
{
    public function run(): void
    {
        Log::info('Starting admin user seeding');

        $adminRole = Role::firstOrCreate(['name' => 'admin']);

        $adminUser = User::firstOrCreate(
            ['email' => 'admin@ags.com'],
            [
                'name' => 'Admin Dummy',
                'password' => Hash::make('password'),
            ]
        );
        $adminUser->roles()->sync([$adminRole->id]);

        Log::info('Admin user created and assigned role', ['email' => $adminUser->email, 'role' => $adminRole->name]);

        Log::info('Admin user seeding completed');
    }
}
