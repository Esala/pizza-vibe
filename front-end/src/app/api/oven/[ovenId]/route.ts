import { NextRequest, NextResponse } from 'next/server';

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ ovenId: string }> }
) {
  try {
    const { ovenId } = await params;
    const url = new URL(request.url);
    const user = url.searchParams.get('user') || 'user';

    const ovenServiceUrl = process.env.OVEN_SERVICE_URL || 'http://localhost:8085';
    const response = await fetch(`${ovenServiceUrl}/ovens/${ovenId}?user=${user}`, {
      method: 'POST',
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to reserve oven' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Error reserving oven:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ ovenId: string }> }
) {
  try {
    const { ovenId } = await params;
    const ovenServiceUrl = process.env.OVEN_SERVICE_URL || 'http://localhost:8085';
    const response = await fetch(`${ovenServiceUrl}/ovens/${ovenId}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to release oven' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Error releasing oven:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
