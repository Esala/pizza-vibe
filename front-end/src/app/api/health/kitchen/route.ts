import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const kitchenServiceUrl = process.env.KITCHEN_SERVICE_URL || 'http://localhost:8081';
    const response = await fetch(`${kitchenServiceUrl}/health`, {
      method: 'GET',
    });

    if (!response.ok) {
      return NextResponse.json(
        { status: 'unhealthy' },
        {
          status: 503,
          headers: {
            'Access-Control-Allow-Origin': '*',
          },
        }
      );
    }

    const data = await response.json();
    return NextResponse.json(data, {
      status: 200,
      headers: {
        'Access-Control-Allow-Origin': '*',
      },
    });
  } catch (error) {
    console.error('Kitchen service health check failed:', error);
    return NextResponse.json(
      { status: 'unhealthy' },
      {
        status: 503,
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
      }
    );
  }
}

export async function OPTIONS() {
  return new NextResponse(null, {
    status: 200,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type',
    },
  });
}
