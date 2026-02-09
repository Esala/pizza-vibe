import { NextRequest, NextResponse } from 'next/server';

interface AddQuantityRequest {
  quantity: number;
}

function validateAddQuantityRequest(body: unknown): body is AddQuantityRequest {
  if (!body || typeof body !== 'object') {
    return false;
  }

  const req = body as Record<string, unknown>;
  return typeof req.quantity === 'number' && req.quantity > 0;
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ item: string }> }
) {
  try {
    const { item } = await params;
    const body = await request.json();

    if (!validateAddQuantityRequest(body)) {
      return NextResponse.json(
        { error: 'Invalid request body' },
        { status: 400 }
      );
    }

    const inventoryServiceUrl = process.env.INVENTORY_SERVICE_URL || 'http://localhost:8084';
    const response = await fetch(`${inventoryServiceUrl}/inventory/${item}/add`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to add quantity' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Error adding quantity:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
