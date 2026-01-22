import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const kitchenServiceUrl = process.env.KITCHEN_SERVICE_URL || 'http://localhost:8081';
    const response = await fetch(`${kitchenServiceUrl}/orders`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to fetch kitchen orders' },
        {
          status: response.status,
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
    console.error('Error fetching kitchen orders:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      {
        status: 500,
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
