import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const deliveryServiceUrl = process.env.DELIVERY_SERVICE_URL || 'http://localhost:8082';
    const response = await fetch(`${deliveryServiceUrl}/health`, {
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
    console.error('Delivery service health check failed:', error);
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
